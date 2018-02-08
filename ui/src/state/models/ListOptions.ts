import { observable } from 'mobx';

type dir = 'asc' | 'desc';

export default class ListOptions {
    @observable per_page: number = 20;
    @observable page: number = 0;
    @observable sort_by: string = '';
    @observable sort_dir: dir = 'asc';

    toJS(): Partial<ListOptions> {
        const { per_page, page, sort_by, sort_dir } = this;
        return { per_page, page, sort_by, sort_dir };
    }
}
