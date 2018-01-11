import * as React from 'react';
import { observer, inject } from "mobx-react";

import RootStore from '../../state/RootStore';
import PaginationModel from '../../state/models/Pagination';
import './Pagination.scss';

interface Props {
    pagination: PaginationModel;
    maxNumbers?: number;
    store?: RootStore;
}

@inject('store')
@observer
export default class Pagination extends React.Component<Props> {
    static defaultProps: Partial<Props> = {
        maxNumbers: 10
    }

    handlePageNumberClick = (pageNumber: number) => {
        this.props.store.images.loadPage(pageNumber);
    }

    renderPageNumber = (number: number) => {
        const disabled = number === (this.props.pagination.currentPage + 1);
        return <a onClick={() => !disabled && this.handlePageNumberClick(number - 1)} key={number} className="pagination-page-number" disabled={disabled}>{number}</a>;
    }

    renderPageNumbers() {
        const { pagination, maxNumbers } = this.props;
        const surroundingPages = pagination.surroundingPages;
        const numSurroundings = surroundingPages.length;

        const addFirstPage = surroundingPages[0] != 0;
        const addLastPage = surroundingPages[numSurroundings - 1] != pagination.totalPages - 1;

        const numberElements = surroundingPages.map(i => this.renderPageNumber(i + 1));

        if (numSurroundings < pagination.maxVisibleNumbers) {
            return numberElements;
        }

        if (addFirstPage) {
            if (surroundingPages[0] == 1) {
                numberElements.unshift(this.renderPageNumber(1));
            } else {
                numberElements.unshift(<span key="dots-1" className="pagination-dots">...</span>);
                numberElements.unshift(this.renderPageNumber(1));
            }
        }

        if (addLastPage) {
            if (surroundingPages[numSurroundings - 1] == pagination.totalPages - 2) {
                numberElements.push(this.renderPageNumber(pagination.totalPages));
            } else {
                numberElements.push(<span key="dots-1" className="pagination-dots">...</span>);
                numberElements.push(this.renderPageNumber(pagination.totalPages));
            }
        }

        return numberElements;
    }

    render() {
        const { pagination } = this.props;
        const classes = ['pagination'].concat([this.props.className || '']).join(' ');
        return (
            <div className={classes}>
                <div className="pagination-items">
                    SHOWING: {pagination.currentItemsRange}
                </div>
                <div className="pagination-pages">
                    PAGES:
                    <div className="pagination-page-numbers">
                        {this.renderPageNumbers()}
                    </div>
                </div>
            </div>
        );
    }
}
