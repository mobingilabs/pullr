import React from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';

export class ModalBoundry extends React.PureComponent {
    static childContextTypes = {
        modalBoundary: PropTypes.object
    }

    getChildContext() {
        return {
            modalBoundary: {
                get: () => this.modalsRef
            }
        };
    }

    render() {
        const classes = ['modalboundary'].concat([
            this.props.className || ''
        ]).join(' ');

        return (
            <div className={classes}>
                {this.props.children}
                <div className="modals" ref={e => this.modalsRef = e }></div>
            </div>
        );
    }
}

export default withRouter(ModalBoundry);