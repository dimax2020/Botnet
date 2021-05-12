export interface Botinfo {
  id: number;
  address: string;
  ping: number;
  CommandExecStatus: string;
  status: string;
  commandOutput?: string;
}
export interface Postobj {
  command: string;
  option: string;

}
export interface BotsInfo extends Array<Botinfo>{}
export interface Postobjs extends Array<Postobj>{}
