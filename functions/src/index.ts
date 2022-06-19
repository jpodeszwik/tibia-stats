import * as functions from "firebase-functions";
import * as dateFormat from "dateformat";
import * as admin from "firebase-admin";
import * as tibiadata from "./tibiadata";

const PELORIA = "Peloria";
const UPDATE_MBEMBERS_COMMANDS = "updateMembersCommands";
const UPDATE_EXPERIENCE_COMMANDS = "updatExperienceCommands";

admin.initializeApp();

const db = admin.firestore();

type ExpData = {
    name: string,
    vocation: string,
    exp: number
};

const saveExpData = async (expData: ExpData[]) => {
  const batchSize = 500;
  const date = dateFormat(new Date(), "yyyy-mm-dd");

  while (expData.length > 0) {
    const batchRecords = expData.splice(0, batchSize);
    console.log(`Inserting batch of size ${batchRecords.length}`);

    const batch = db.batch();
    batchRecords.forEach(
        (val) => batch.set(db.doc(`date/${date}/char/${val.name}`), val));
    await batch.commit();
  }
};

const updateExp = async (world: string) => {
  const expData = await tibiadata.fetchAllProfessionsHighscore(world, tibiadata.Highscore.EXP);

  await saveExpData(expData.map((val) => ({name: val.name, vocation: val.vocation, exp: val.value})));
};

const updateGuildMembers = async (world: string) => {
  const guilds = await tibiadata.fetchGuilds(world);

  await Promise.all(guilds.map((guild) => tibiadata.fetchGuild(guild.name)
      .then((guildDetails) => saveGuildMembers(guildDetails.name, guildDetails.members.map((member) => member.name)))
  ));
};

const saveGuildMembers = (guildName: string, members: string[]) => {
  const date = dateFormat(new Date(), "yyyy-mm-dd");
  return db.doc(`date/${date}/guild/${guildName}`).set({members: members});
};

const scheduleMembersUpdate = (world: string) =>
  db.collection(UPDATE_MBEMBERS_COMMANDS).add({
    date: dateFormat(new Date(), "yyyy-mm-dd HH:MM"),
    world: world,
  });

const scheduleExperienceUpdate = (world: string) =>
  db.collection(UPDATE_EXPERIENCE_COMMANDS).add({
    date: dateFormat(new Date(), "yyyy-mm-dd HH:MM"),
    world: world,
  });


exports.forceUpdateMembers = functions.https.onRequest((req, res) => {
  scheduleMembersUpdate(PELORIA)
      .finally(() => res.end());
});

exports.updateMembers = functions.pubsub
    .schedule("10 10 * * * ")
    .timeZone("Europe/Warsaw")
    .onRun(() => scheduleMembersUpdate(PELORIA));

exports.forceUpdateExperience = functions.https.onRequest((req, res) => {
  scheduleExperienceUpdate(PELORIA)
      .finally(() => res.end());
});

exports.updateExperience = functions.pubsub
    .schedule("10 10 * * * ")
    .timeZone("Europe/Warsaw")
    .onRun(() => scheduleExperienceUpdate(PELORIA));

exports.onUpdateMembers = functions.firestore
    .document(`${UPDATE_MBEMBERS_COMMANDS}/{commandId}`)
    .onCreate((snap) => updateGuildMembers(snap.data().world));

exports.onUpdateExperience = functions.firestore
    .document(`${UPDATE_EXPERIENCE_COMMANDS}/{commandId}`)
    .onCreate((snap) => updateExp(snap.data().world));

