export class IbcStatisticsResDto {
  statistics_name: string;
  count: string;

  constructor(value) {
    const { statistics_name, count } = value;
    this.statistics_name = statistics_name;
    this.count = count;
  }

  static bundleData(value: any): IbcStatisticsResDto[] {
    const datas: IbcStatisticsResDto[] = value.map((item: any) => {
      return new IbcStatisticsResDto(item);
    });
    return datas;
  }
}
