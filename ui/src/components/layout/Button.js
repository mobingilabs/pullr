import React from 'react';
import PropTypes from 'prop-types';
import { Icon } from 'react-fa';

import './Button.scss';

export default class Button extends React.PureComponent {
    static propTypes = {
        text: PropTypes.node,
        icon: PropTypes.string,
        onClick: PropTypes.func.isRequired,
        disabled: PropTypes.bool,
        className: PropTypes.string,
        size: PropTypes.oneOf(['big', 'small', 'medium']),
        secondary: PropTypes.bool
    }

    static defaultProps = {
        disabled: false,
        size: 'medium'
    }

    render () {
        const classes = ['button', this.props.size].concat([
            this.props.className ? this.props.className : '',
            this.props.text      ? '' : 'button-icon-only',
            this.props.secondary ? 'button-secondary' : ''
        ]).join(' ');

        return (
            <button disabled={ this.props.disabled } className={ classes } onClick={this.props.onClick}>
                { this.props.icon ? <Icon name={ this.props.icon } /> : null }
                { this.props.text }
            </button>
        )
    }
}
