import { observable, computed, action } from "mobx";

export class PageOutOfRangeException extends Error {}

export default class Pagination {
    @observable currentPage: number = 0;
    @observable itemsPerPage: number = 10;
    @observable totalItems: number = 0;
    @observable maxVisibleNumbers: number = 10;

    constructor(data?: Partial<Pagination>) {
        if (data) {
            this.currentPage = data.currentPage || this.currentPage;
            this.itemsPerPage = data.itemsPerPage || this.itemsPerPage;
            this.totalItems = data.totalItems || this.totalItems;
        }
    }

    @computed get totalPages(): number {
        return Math.ceil(this.totalItems / this.itemsPerPage);
    }

    @computed get currentItemsRange(): string {
        return `${this.currentPage * this.itemsPerPage + 1} - ${Math.min(this.currentPage + this.itemsPerPage + 1, Math.max(this.totalItems, 1))}`;
    }

    @computed get surroundingPages(): Array<number> {
        if (this.totalPages == 0 || this.totalPages == 1) {
            return [0];
        }

        const halfVisibleNumbers = Math.ceil(this.maxVisibleNumbers / 2.0);
        const end = Math.min(this.totalPages - 1, this.currentPage + halfVisibleNumbers);
        const leftover = this.maxVisibleNumbers - (end - this.currentPage);
        const start = Math.max(0, this.currentPage - halfVisibleNumbers - leftover);
        console.log(`Calculating surrounding pages ${start}-${end}`);
        return Array<number>(end - start).keys().map((i: number) => i + start);
    }

    @action gotoPage(pageIndex: number) {
        if (pageIndex < 0 || pageIndex >= this.totalPages) {
            throw new PageOutOfRangeException();
        }

        this.currentPage = pageIndex;
    }
}
