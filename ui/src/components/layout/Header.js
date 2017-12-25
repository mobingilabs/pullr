import React from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { Icon } from 'react-fa';

import Icons from './Icons';
import Button from './Button';
import './Header.scss';

export class Header extends React.PureComponent {
    static propTypes = {
        title: PropTypes.string.isRequired,
        subTitle: PropTypes.node,
        back: PropTypes.string,
        actions: PropTypes.arrayOf(PropTypes.shape({
            text: PropTypes.string.isRequired,
            handler: PropTypes.func.isRequired,
            icon: PropTypes.string,
            disabled: PropTypes.bool
        }))
    }

    static defaultProps = {
        subTitle: '',
        actions: []
    }

    render () {
        const classes = ['screen-header'].concat([
            this.props.back ? 'has-breadcrumb' : ''
        ]).join(' ');

        return (
            <div className={ classes }>
                <div className="titlewrapper">
                    { this.props.back &&
                            <a className="header-backbutton" onClick={ this.props.history.goBack }>
                                <Icon name={ Icons.back } /> BACK
                            </a>
                    }
                    <h1 className="title">{ this.props.title }</h1>
                    <div className="subtitle">{ this.props.subTitle }</div>
                </div>
                { this.props.actions &&
                    <div className="actions">
                        { this.props.actions.map((action, i) => 
                            <Button key={i} 
                                className="button-header" 
                                text={action.text} 
                                onClick={action.handler} 
                                icon={action.icon} 
                                disabled={action.disabled} 
                                />
                        )}
                    </div>
                }
            </div>
        );
    }
}

export default withRouter(Header);
