import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Route, Link } from 'react-router-dom';
import { Icon } from 'react-fa';

import Screen from '../layout/Screen';
import Header from '../layout/Header';
import Button from '../layout/Button';
import Icons from '../layout/Icons';
import Pagination from '../widgets/Pagination';

import * as actions from '../../state/history/actions';

export class HistoryScreen extends React.PureComponent {
    render () {
        const loadMoreVisible = !this.props.loadInProgress && this.props.thereIsMore;
        const loadingSpinnerVisible = this.props.loadInProgress;
        return (
            <Screen>
                <Header title="BUILD HISTORY" subTitle={`Last build on ${new Date()}`} />
                <div className="scroll-horizontal">
                    <table className="wide big-shadow">
                        <thead>
                            <tr>
                                <th>IMAGE NAME</th>
                                <th>BUILD TAG</th>
                                <th>DATE</th>
                                <th>STATUS</th>
                            </tr>
                        </thead>
                        <tbody>
                            {this.props.lastBuilds.map(build =>
                                <tr key={build.imageName}>
                                    <td>{build.imageName}</td>
                                    <td>{build.date.toString()}</td>
                                    <td>{build.tag}</td>
                                    <td>{build.status}</td>
                                </tr>
                            )}
                            { (loadMoreVisible || loadingSpinnerVisible) &&
                                <tr className="load-more">
                                    <td colSpan="4">
                                        { loadMoreVisible && 
                                        <Button onClick={ this.props.onLoadMore } text="LOAD MORE" /> }
                                        { loadingSpinnerVisible &&
                                        <Icon spin name={ Icons.loadingSpinner } /> }
                                    </td>
                                </tr>
                            }
                        </tbody>
                    </table>
                </div>
            </Screen>
        );
    }
}

const mapStateToProps = (state) => {
    return state.history;
};

const mapDispatchToProps = (dispatch) => {
    return {
        onGotoPage: (pageNumber) => dispatch(actions.gotoPage(pageNumber)),
        onLoadMore: () => dispatch(actions.loadMore())
    };
};

export default connect(mapStateToProps, mapDispatchToProps)(HistoryScreen);
