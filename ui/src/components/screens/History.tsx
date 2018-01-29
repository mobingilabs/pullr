import * as React from 'react';
import { Icon } from 'react-fa';
import { observer, inject } from "mobx-react";

import Screen from '../layout/Screen';
import Header from '../layout/Header';
import Button from '../layout/Button';
import Icons from '../layout/Icons';
import Pagination from '../widgets/Pagination';

@observer
export default class HistoryScreen extends React.Component {
    render () {
        const loadMoreVisible = true;
        const loadingSpinnerVisible = false;
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
                            {/* {this.props.store.lastBuilds.map(build =>
                                <tr key={build.imageName}>
                                    <td>{build.imageName}</td>
                                    <td>{build.date.toString()}</td>
                                    <td>{build.tag}</td>
                                    <td>{build.status}</td>
                                </tr>
                            )} */}
                            { (loadMoreVisible || loadingSpinnerVisible) &&
                                <tr className="load-more">
                                    <td colSpan={4}>
                                        {/* { loadMoreVisible && 
                                        <Button onClick={this.props.onLoadMore} text="LOAD MORE" /> 
                                        }
                                        { loadingSpinnerVisible &&
                                        <Icon spin name={ Icons.loadingSpinner } /> } */}
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
