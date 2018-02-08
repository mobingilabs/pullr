import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { observable, action } from 'mobx';
import { observer } from 'mobx-react';
import { Icon } from 'react-fa';

import { ListOptions } from '../../libs/api/ApiClient';
import Icons from '../layout/Icons';

import './TableActionsMenu.scss';

interface Props {
    listOptions?: Partial<ListOptions>;
    portal: HTMLElement;
}

@observer
export default class TableActionsMenu extends React.Component<Props> {
    @observable menuShown: boolean = false;
    editorEl: TableActionsEditor;

    @action.bound
    toggleMenu() {
        this.menuShown = !this.menuShown;
        if (this.menuShown) {
            document.addEventListener('click', this.handleDocumentClick);
        } else {
            document.removeEventListener('click', this.handleDocumentClick);
        }
    }

    @action.bound
    handleDocumentClick(e: MouseEvent) {
        const el = ReactDOM.findDOMNode(this.editorEl);
        const contains = el.contains(e.target as Node);
        if (el != e.target && !contains) {
            this.toggleMenu();
        }
    }

    render() {
        let menu = null;
        if (this.menuShown) {
            menu = ReactDOM.createPortal(<TableActionsEditor key="menu" ref={e => this.editorEl = e} options={this.props.listOptions} />, this.props.portal);
        }

        return [
            <div className="th table-actions-menu" onClick={this.toggleMenu} key="menu-icon">
                <Icon name={Icons.menu} />
            </div>,
            menu
        ]
    }
}

interface EditorProps {
    options: Partial<ListOptions>;
}

const defined = (val: any): boolean => {
    return val !== undefined && val !== null;
}

@observer
export class TableActionsEditor extends React.Component<EditorProps> {
    render() {
        // TODO: Implement editor
        return (
            <div className="table-actions-editor">

            </div>
        )
    }
}
