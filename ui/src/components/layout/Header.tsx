import * as React from 'react';
import { Icon } from 'react-fa';
import { withRouter, RouteComponentProps } from 'react-router-dom';

import Icons from './Icons';
import Button from './Button';
import './Header.scss';

interface Action {
    text: string;
    handler: (_: any) => any;
    icon?: string;
    disabled?: boolean;
}

interface Props {
    title: string;
    subTitle?: string;
    back?: boolean;
    actions?: Array<Action>
}

@withRouter
export default class Header extends React.PureComponent<Props & Partial<RouteComponentProps<{}>> {
    public static defaultProps: Partial<Props> = {
        subTitle: '',
        actions: [],
        back: false
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
