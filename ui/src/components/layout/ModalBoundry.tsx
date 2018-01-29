import * as React from 'react';
import * as PropTypes from 'prop-types';
import { withRouter } from "react-router-dom";

interface Props {
    modalRoot: HTMLElement;
}

@withRouter
export default class ModalBoundary extends React.PureComponent<Props> {
    modalsRef: HTMLDivElement;

    static childContextTypes = {
        modalRoot: PropTypes.any
    }

    getChildContext() {
        return {
            modalRoot: this.props.modalRoot
        };
    }

    render() {
        return this.props.children;
    }
}
