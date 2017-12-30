import React from 'react';
import PropTypes from 'prop-types';

import Screen from '../layout/Screen';
import Header from '../layout/Header';

export default class HistoryScreen extends React.PureComponent {
    render () {
        return (
            <Screen>
                <Header title="BUILD HISTORY" subTitle={`Last build on ${new Date()}`} />
                <div className="scroll-horizontal">
                    <table className="wide">
                        <thead>
                            <tr>
                                <th>IMAGE NAME</th>
                                <th>BUILD TAG</th>
                                <th>DATE</th>
                                <th>STATUS</th>
                            </tr>
                        </thead>
                        <tbody>
                            {builds.map(build =>
                                <tr key={images[key].name}>
                                    <td>
                                        <Link className="table-link" to={`/images/${images[key].name}`}><Icon name={Icons.images} /> {images[key].name}</Link>
                                    </td>
                                    <td>{images[key].provider}</td>
                                    <td>{images[key].organisation}/{images[key].repository}</td>
                                    <td>{images[key].builds.map(build => build.tag || build.name).join(', ')}</td>
                                    <td width="100"><Button icon={Icons.edit} onClick={this.handleEditImage.bind(null, key)} /></td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>
                <Pagination
                    className="big-shadow"
                    currentPage={this.props.currentPage}
                    totalPages={this.props.totalPages}
                    onChangePage={this.props.onShowPage}
                    itemPerPage={10}
                    totalItems={this.props.totalImages} />

                <Route path="/images/:imageName" component={ImageDetailModal} />
            </Screen>
        )
    }
}

