import React from 'react';
import PropTypes from 'prop-types';
import { Icon } from 'react-fa';

import Button from '../layout/Button';
import './Card.scss';

export default class Card extends React.PureComponent {
    static propTypes = {
        icon: PropTypes.string.isRequired,
        title: PropTypes.string.isRequired,
        background: PropTypes.string,
        disabled: PropTypes.bool,
        dark: PropTypes.bool
    }

    static defaultProps = {
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
