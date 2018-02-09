import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { observable, action } from 'mobx';
import { observer } from 'mobx-react';
import * as Popover from 'react-popover';
import { Icon } from 'react-fa';

import ListOptions from '../../state/models/ListOptions';
import TableOptions from '../../state/models/TableOptions';
import Icons from '../layout/Icons';

import './TableActionsMenu.scss';

interface Props {
    tableOptions?: Partial<TableOptions>;
    listOptions?: Partial<ListOptions>;
    sortableFields?: { [field: string]: string };
    portal: HTMLElement;
}

@observer
export default class TableActionsMenu extends React.Component<Props> {
    @observable menuShown: boolean = false;

    @action.bound
    toggleMenu() {
        this.menuShown = !this.menuShown;
    }

    @action.bound
    closeMenu(e: Event) {
        this.menuShown = false;
    }

    render() {
        const editor = (
            <TableActionsEditor
                sortableFields={this.props.sortableFields}
                tableOptions={this.props.tableOptions}
                listOptions={this.props.listOptions} />
        );

        const tableFiltered = this.props.tableOptions && this.props.tableOptions.query.length > 0;
        const menuClasses = ['table-actions-menu'].concat([
            tableFiltered ? 'table-actions-menu-filtered' : ''
        ]).join(' ');

        const menuTitle = tableFiltered ? 'Results are filtered' : null;

        return (
            <Popover className="table-actions-popover" isOpen={this.menuShown} preferPlace="below" appendTarget={this.props.portal} body={editor} onOuterAction={this.closeMenu}>
                <div className={menuClasses} title={menuTitle} onClick={this.toggleMenu} key="menu-icon">
                    <Icon name={Icons.menu} />
                </div>
            </Popover>
        )
    }
}

interface EditorProps {
    tableOptions?: Partial<TableOptions>;
    listOptions?: Partial<ListOptions>;
    sortableFields?: { [field: string]: string };
}

const defined = (val: any): boolean => {
    return val !== undefined && val !== null;
}

@observer
export class TableActionsEditor extends React.Component<EditorProps> {
    bindTableOpts = (field: string) => {
        return (e: any) => {
            const val = e.target.value;
            (this.props.tableOptions as any)[field] = val;
        }
    }

    bindListOpts = (field: string) => {
        return (e: any) => {
            const val = e.target.value;
            (this.props.listOptions as any)[field] = val;
        }
    }
    render() {
        return (
            <div className="table-actions-editor">
                <div className="form">
                    {this.props.tableOptions &&
                        <div className="entry">
                            <label htmlFor="table-actions-query">Filter</label>
                            <input id="table-actions-query" type="text" value={this.props.tableOptions.query} onChange={this.bindTableOpts('query')} />
                        </div>
                    }
                    <div className="entry">
                        <label htmlFor="table-actions-per-page">Items Per Page</label>
                        <select id="table-actions-per-page" value={this.props.listOptions.per_page} onChange={this.bindListOpts('per_page')}>
                            <option value={20}>20</option>
                            <option value={50}>50</option>
                            <option value={100}>100</option>
                        </select>
                    </div>
                    {this.props.sortableFields &&
                        <div className="entry">
                            <label htmlFor="table-actions-sort-by">Sort By</label>
                            <select id="table-actions-sort-by" value={this.props.listOptions.sort_by} onChange={this.bindListOpts('sort_by')}>
                                <option key="null" value={null}>Default</option>
                                {Object.keys(this.props.sortableFields).map(field =>
                                    <option key={field} value={field}>{this.props.sortableFields[field]}</option>
                                )}
                            </select>
                        </div>
                    }
                    {this.props.sortableFields && this.props.listOptions.sort_by != "" &&
                        <div className="entry">
                            <label htmlFor="table-actions-sort-dir">Sort Direction</label>
                            <select id="table-actions-sort-dir" value={this.props.listOptions.sort_dir} onChange={this.bindListOpts('sort_dir')}>
                                <option value={null}>Default</option>
                                <option value="asc">Ascending</option>
                                <option value="desc">Descending</option>
                            </select>
                        </div>
                    }
                </div>
            </div>
        )
    }
}
