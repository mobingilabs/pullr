import * as React from 'react';
import { NavLink, withRouter } from 'react-router-dom';
import { Icon } from 'react-fa';

interface Props {
    path?: string;
    onClick?: (e: any) => any,
    icon: string;
    text: string;
}

@withRouter
export default class MenuItem extends React.PureComponent<Props> {
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

        return <li>{link}</li>;
    }
}
