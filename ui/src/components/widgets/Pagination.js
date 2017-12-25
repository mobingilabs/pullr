import React from 'react';
import PropTypes from 'prop-types';

import './Pagination.scss';

export default class Pagination extends React.PureComponent {
    static propTypes = {
        currentPage: PropTypes.number.isRequired,
        totalPages: PropTypes.number.isRequired,
        totalItems: PropTypes.number.isRequired,
        onChangePage: PropTypes.func.isRequired,
        itemPerPage: PropTypes.number.isRequired,
        maxPages: PropTypes.number
    }

    static defaultProps = {
        maxPages: 10
    }

    renderShowedItemInfo () {
        const start = (this.props.itemPerPage * this.props.currentPage) + 1;
        const stop = Math.min(((start - 1) + this.props.itemPerPage), this.props.totalItems);
        return `${start}-${stop} / ${this.props.totalItems}`;
    }

    renderPageNumber = (number) => {
        const disabled = number === (this.props.currentPage + 1);
        return <a onClick={() => !disabled && this.props.onChangePage(number - 1)} key={number} className="pagination-page-number" disabled={disabled}>{ number }</a>;
    }

    renderPageNumbers () {
        const showFirstPage = this.props.currentPage > Math.floor(this.props.maxPages / 2.0);
        const showLastPage = this.props.totalPages - this.props.currentPage > Math.floor(this.props.maxPages / 2.0);

        let pages = [];
        let start = Math.max(this.props.currentPage - Math.floor(this.props.maxPages / 2.0), 0);
        let stop = Math.min(start + this.props.maxPages, this.props.totalPages);
        for (let i = start; i < stop; i++) {
            pages.push(this.renderPageNumber(i + 1));
        }

        if (showFirstPage) {
            if (start != 1) {
                pages.unshift(<span key="dots-1" className="pagination-dots">...</span>);
            }

            pages.unshift(this.renderPageNumber(1));
        }


        if (showLastPage) {
            if (stop != this.props.totalPages - 1) {
                pages.push(<span key="dots-2" className="pagination-dots">...</span>);
            }
            pages.push(this.renderPageNumber(this.props.totalPages));
        }

        return pages;
    }

    render () {
        const classes = ['pagination'].concat([this.props.className || '']).join(' ');
        return (
            <div className={ classes }>
                <div className="pagination-items">
                    SHOWING: { this.renderShowedItemInfo() }
                </div>
                <div className="pagination-pages">
                    PAGES:
                    <div className="pagination-page-numbers">
                        { this.renderPageNumbers() }
                    </div>
                </div>
            </div>
        );
    }
}
