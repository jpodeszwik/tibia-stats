import * as functions from "firebase-functions";
import fetch from "node-fetch";
import * as dateFormat from "dateformat";
import * as admin from "firebase-admin";

const WORLD = "Peloria";
const API = "https://api.tibiadata.com/v3";

admin.initializeApp();

const db = admin.firestore();

type ExpData = {
    name: string,
    vocation: string,
    exp: number
};

const fetchAllProfessions: (world: string) => Promise<ExpData[]> =
  async (world: string) => {
    const professions = ["knight", "paladin", "druid", "sorcerer"]
        .map((profession) => fetchProfession(world, profession));

    const results = await Promise.all(professions);
    return results.reduce((a, b) => a.concat(b));
  };

const fetchProfession: (world: string, profession: string) =>
  Promise<ExpData[]> =
    async (world: string, profession: string) => {
      const url = `${API}/highscores/${world}/exp/${profession}`;
      console.log(url);
      return fetch(url).then((res) => res.json())
          .then((res) => res.highscores.highscore_list)
          .then((rows: any[]) => rows.map((row) =>
            ({name: row.name, vocation: row.vocation, exp: row.value})));
    };

const updateExp = async () => {
  const batchSize = 500;
  const date = dateFormat(new Date(), "yyyy-mm-dd");
  const res = await fetchAllProfessions(WORLD);

  while (res.length > 0) {
    const batchRecords = res.splice(0, batchSize);
    console.log(`Inserting batch of size ${batchRecords.length}`);

    const batch = db.batch();
    batchRecords.forEach(
        (val) => batch.set(db.doc(`date/${date}/char/${val.name}`), val));
    await batch.commit();
  }

  await db.doc("metadata/lastScan")
      .set({time: dateFormat(new Date(), "yyyy-mm-dd HH:MM")});
};

const fetchGuilds: (world: string) => Promise<string[]> =
  async (world: string) => {
    const url = `${API}/guilds/${world}`;
    return fetch(url)
        .then((res) => res.json())
        .then((res) => res.guilds.active
            .map((guild: any) => guild.name)
        );
  };

const fetchGuildMembers: (guildName: string) => Promise<string[]> =
  async (guildName: string) => {
    const url = `${API}/guild/${guildName}`;
    return fetch(url)
        .then((res) => res.json())
        .then((res) => {
          if (res.guilds.guild.error) {
            throw res.guilds.guild.error;
          }
          return res.guilds.guild;
        })
        .then(parseGuild);
  };

const parseGuild: (guild: any) => string[] = (guild: any) => guild.members
    .map((member: any) => member.name);

const updateGuildMembers = async () => {
  const date = dateFormat(new Date(), "yyyy-mm-dd");
  const guilds = await fetchGuilds(WORLD);

  for (let i = 0; i < guilds.length; i++) {
    const guild = guilds[i];
    await fetchGuildMembers(guild)
        .then((members) =>
          db.doc(`date/${date}/guild/${guild}`).set(
              {members: members}))
        .catch((err) =>
          console.error(`Failed to fetch guild ${guild}: ${err}`));
  }
};

exports.forceUpdateMembers = functions.https.onRequest((req, res) => {
  updateGuildMembers().finally(() => res.end());
});

exports.updateMembers = functions.pubsub
    .schedule("10 10 * * * ")
    .timeZone("Europe/Warsaw")
    .onRun((context) => updateGuildMembers());

exports.forceUpdateExperience = functions.https.onRequest((req, res) => {
  updateExp().finally(() => res.end());
});

exports.updateExperience = functions.pubsub
    .schedule("10 10 * * * ")
    .timeZone("Europe/Warsaw")
    .onRun((context) => updateExp());
