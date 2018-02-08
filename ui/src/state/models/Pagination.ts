import { observable, computed, action } from "mobx";

export class PageOutOfRangeException extends Error { }

export default class Pagination {
    @observable current: number = 0;
    @observable next: number = 0;
    @observable last: number = 0;
    @observable per_page: number = 10;
    @observable total: number = 0;
    @observable maxVisibleNumbers: number = 10;

    constructor(data?: Partial<Pagination>) {
        if (data) {
            this.current = data.current || this.current;
            this.per_page = data.per_page || this.per_page;
            this.total = data.total || this.total;
            this.next = data.next || this.next;
            this.last = data.last || this.last;
            this.maxVisibleNumbers = data.maxVisibleNumbers || this.maxVisibleNumbers;
        }
    }

    @computed get currentItemsRange(): string {
        return `${this.current * this.per_page + 1} - ${Math.min(this.current + this.per_page + 1, Math.max(this.total, 1))}`;
    }

    @computed get surroundingPages(): Array<number> {
        if (this.last == 0 || this.last == 1) {
            return [0];
        }

        const halfVisibleNumbers = Math.ceil(this.maxVisibleNumbers / 2.0);
        const end = Math.min(this.last - 1, this.current + halfVisibleNumbers);
        const leftover = this.maxVisibleNumbers - (end - this.current);
        const start = Math.max(0, this.current - halfVisibleNumbers - leftover);
        console.log(`Calculating surrounding pages ${start}-${end}`);
        return Array<number>(end - start).keys().map((i: number) => i + start);
    }

    @action gotoPage(pageIndex: number) {
        if (pageIndex < 0 || pageIndex >= this.last) {
            throw new PageOutOfRangeException();
        }

        this.current = pageIndex;
    }
}
