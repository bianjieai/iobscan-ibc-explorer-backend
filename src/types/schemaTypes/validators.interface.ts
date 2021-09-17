export interface IValidatorsStruct {
  name?:string,
  pubkey?:string,
  power?:string,
  operator?:string,
  jailed?:boolean | string,
  details?:string;
  hash?: string;
}

