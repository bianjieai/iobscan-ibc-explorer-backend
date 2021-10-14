export class ConfigResDto {
  iobscan: string

  constructor(value) {
    const { iobscan } = value;
    this.iobscan = iobscan
  }

}
