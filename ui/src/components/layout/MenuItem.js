import React from 'react';
import PropTypes from 'prop-types';
import { Icon } from 'react-fa';
import { NavLink, withRouter } from 'react-router-dom';

export class MenuItem extends React.PureComponent {
    static propTypes = {
        path: PropTypes.string,
        onClick: PropTypes.func,
        icon: PropTypes.string.isRequired,
        text: PropTypes.string.isRequired
    }

    render () {
        let link;
        if (this.props.onClick) {
            link = (
                <a onClick={ this.props.onClick }>
                    <Icon name={this.props.icon} />
                    { this.props.text }
                </a>
            );
        } else {
            link = (
                <NavLink to={ this.props.path } activeClassName="active">
                    <Icon name={ this.props.icon } />
                    { this.props.text }
                </NavLink>
            );
        }

        return <li>{ link }</li>;
    }
}

export default withRouter(MenuItem);
