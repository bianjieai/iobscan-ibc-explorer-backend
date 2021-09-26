export class IbcStatisticsResDto {
  statistics_name: string;
  count: string;
  create_at: string;
  update_at: string;
  constructor(value) {
    const { statistics_name, count, create_at, update_at } = value;
    this.statistics_name = statistics_name;
    this.count = count;
    this.create_at = create_at;
    this.update_at = update_at;
  }
}
