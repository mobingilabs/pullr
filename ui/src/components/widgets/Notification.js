import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Icon } from 'react-fa';

import Icons from '../layout/Icons';
import * as notificationActions from '../../state/notifications/actions';
import * as notificationSelectors from '../../state/notifications/selectors';

import './Notification.scss';

export class Notification extends React.PureComponent {
    static propTypes = {
        id: PropTypes.string.isRequired
    }

    close = () => {
        this.props.onClose(this.props.id);
    }

    render() {
        const { notification } = this.props;
        if (!notification) {
            return null;
        }

        return (
            <div className={`notification ${notification.type}`}>
                <button onClick={this.close}><Icon name={Icons.close} /></button>
                {notification.message}
            </div>
        );
    }
}

const mapStateToProps = (state, ownProps) => {
    const notification = notificationSelectors.getById(state, ownProps.id);
    return { notification };
}

const mapDispatchToProps = (dispatch) => {
    return {
        onClose: (id) => dispatch(notificationActions.remove(id))
    };
};

export default connect(mapStateToProps, mapDispatchToProps)(Notification);
