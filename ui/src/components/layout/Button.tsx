import * as React from 'react';
import { Icon } from 'react-fa';

import './Button.scss';

interface Props {
    aslink?: boolean;
    href?: string;
    popup?: boolean;
    text?: string;
    icon?: string;
    onClick: (e: any) => any;
    disabled?: boolean;
    className?: string;
    size?: 'big' | 'small' | 'medium';
    secondary?: boolean;
}

export default class Button extends React.PureComponent<Props> {
    static defaultProps: Partial<Props> = {
        disabled: false,
        size: 'medium',
        aslink: false,
        href: '',
        popup: false;
    }

    render() {
        const classes = ['button', this.props.size].concat([
            this.props.className ? this.props.className : '',
            this.props.text ? '' : 'button-icon-only',
            this.props.secondary ? 'button-secondary' : ''
        ]).join(' ');

        if (this.props.aslink) {
            return (
                <a disabled={this.props.disabled} target={this.props.popup ? '_blank' : '_self'} onClick={this.props.onClick} className={classes} href={this.props.href}>
                    {this.props.icon ? <Icon name={this.props.icon} /> : null}
                    {this.props.text}
                </a>
            );
        }

        return (
            <button disabled={this.props.disabled} className={classes} onClick={this.props.onClick}>
                {this.props.icon ? <Icon name={this.props.icon} /> : null}
                {this.props.text}
            </button>
        );
    }
}
