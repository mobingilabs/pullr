import * as React from 'react';
import { Icon } from 'react-fa';

import Button from '../layout/Button';
import './Card.scss';

interface Props {
    icon: string;
    title: string;
    background?: string;
    disabled?: boolean;
    dark?: boolean;
}

export default class Card extends React.PureComponent<Props> {
    static defaultProps: Partial<Props> = {
        disabled: false,
        dark: false
    }

    render () {
        const classes = ['card'].concat(
            this.props.disabled ? ['disabled'] : [],
            this.props.dark ? ['dark'] : []
        ).join(' ');

        return (
            <div className={ classes }>
                <div className="card-visual" style={{ background: this.props.background }}>
                    <Icon name={ this.props.icon } />
                    <span>{ this.props.title }</span>
                </div>
                <div className="card-content">
                    { this.props.children }
                </div>
            </div>
        )
    }
}
