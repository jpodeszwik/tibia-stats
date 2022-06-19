import fetch from "node-fetch";

const API = "https://api.tibiadata.com/v3";

export type HighscoreData = {
    name: string,
    vocation: string,
    value: number,
};

export enum Highscore {
    EXP = "exp",
}

export enum Profession {
    KNIGHT = "knight",
    PALADIN = "paladin",
    DRUID = "druid",
    SORCERER = "sorcerer",
}

export const fetchProfessionHighscore =
(world: string, highscore: Highscore, profession: Profession): Promise<HighscoreData[]> => {
  const url = `${API}/highscores/${world}/${highscore}/${profession}`;
  console.log(url);
  return fetch(url).then((res) => res.json())
      .then((res) => res.highscores.highscore_list)
      .then((rows: any[]) => rows.map((row) =>
        ({name: row.name, vocation: row.vocation, value: row.value})));
};

export const fetchAllProfessionsHighscore = (world: string, highscore: Highscore): Promise<HighscoreData[]> => {
  const professions = [Profession.KNIGHT, Profession.PALADIN, Profession.DRUID, Profession.SORCERER]
      .map((profession) => fetchProfessionHighscore(world, highscore, profession));

  return Promise.all(professions)
      .then((results) => results.reduce((a, b) => a.concat(b)));
};

export type Guild = {
    name: string,
};

export type GuildDetails = {
    name: string,
    members: GuildMember[],
};

export type GuildMember = {
    name: string,
};

export const fetchGuilds = (world: string): Promise<Guild[]> => {
  const url = `${API}/guilds/${world}`;
  return fetch(url)
      .then((res) => res.json())
      .then((res) => res.guilds.active
          .map((guild: any) => ({name: guild.name}))
      );
};

export const fetchGuild = (guildName: string): Promise<GuildDetails> => {
  const url = `${API}/guild/${guildName}`;
  return fetch(url)
      .then((res) => res.json())
      .then((res) => res.guilds.guild);
};
