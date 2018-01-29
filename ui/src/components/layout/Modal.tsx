import * as React from 'react';
import * as ReactDOM from 'react-dom';
import * as PropTypes from 'prop-types';

import { Icon } from 'react-fa';
import Icons from './Icons';

import './Modal.scss';
import { KeyboardEvent } from 'react';

interface HeaderProps {
    title?: string;
    subTitle?: string;
    onClose?: (_: any) => any;
    noClose?: boolean;
}

export class ModalHeader extends React.PureComponent<HeaderProps> {
    static defaultProps: Partial<HeaderProps> = {
        noClose: false
    }

    render() {
        const hasTitle = this.props.title || this.props.subTitle;

        return (
            <div className="modal-header">
                <div className="modal-header flex flex-h">
                    {hasTitle &&
                        <div className="modal-title-wrapper flex-grow">
                            {this.props.subTitle &&
                                <div className="modal-subtitle">{this.props.subTitle}</div>
                            }
                            {this.props.title &&
                                <div className="modal-title">{this.props.title}</div>
                            }
                        </div>
                    }
                    {!this.props.noClose &&
                        <div className="flex-shrink flex-align-bottom">
                            <button className="modal-close-btn" onClick={this.props.onClose}>
                                <Icon name={Icons.close} />
                            </button>
                        </div>
                    }
                </div>
            </div>
        )
    }
}

export class ModalContent extends React.PureComponent {
    render() {
        return <div className="modal-content">{ this.props.children }</div>;
    }
}

export class ModalActions extends React.PureComponent {
    render() {
        return <div className="modal-actions">{ this.props.children }</div>;
    }
}

interface ModalProps {
    onClose: () => any;
}

export class Modal extends React.PureComponent<ModalProps> {
    static contextTypes = {
        modalRoot: PropTypes.any
    }

    componentDidMount() {
        document.addEventListener('keyup', this.onKeyUp);
    }

    componentWillUnmount() {
        document.removeEventListener('keyup', this.onKeyUp);
    }

    onKeyUp = (e: any) => {
        if (e.key === 'Escape') {
            this.props.onClose();
        }
    }

    render() {
        const modalElement = (
            <div className="modal-backdrop">
                <div className="modal flex flex-v">
                    {this.props.children}
                </div>
            </div>
        );

        return ReactDOM.createPortal(modalElement, document.getElementById('modal-root'));
    }
}
