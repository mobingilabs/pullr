import React from 'react';
import { Route, Link, withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { Icon } from 'react-fa';
import moment from 'moment';

import Screen from '../layout/Screen';
import Header from '../layout/Header';
import Icons from '../layout/Icons';
import Pagination from '../widgets/Pagination';
import ImageDetailModal from './ImageDetailModal';

import * as ImagesActions from '../../state/images/actions';

export class ImagesScreen extends React.PureComponent {
    constructor (props) {
        super(props);

        this.actions = [
            { text: 'Add Image', icon: 'plus', handler: this.handleAddImage }
        ];
    }

    handleAddImage = () => {
        this.props.history.push('/images/add');
    }

    handleEditImage = (row) => {
        console.log(`Edit image named: ${row['name']}`);
    }

    render () {
        return (
            <Screen>
                <Header title="IMAGES" subTitle={`${this.props.totalImages} images found...`} actions={ this.actions } />
                <div className="scroll-horizontal">
                    <table className="wide">
                        <thead>
                            <tr>
                                <th>IMAGE NAME</th>
                                <th>COMMIT</th>
                                <th>TAG</th>
                                <th>LAST BUILD</th>
                            </tr>
                        </thead>
                        <tbody>
                            { this.props.images.map(image => 
                                <tr key={ image.name }>
                                    <td>
                                        <Link className="table-link" to={ `/images/${image.name}` }><Icon name={ Icons.images }/> { image.name }</Link>
                                    </td>
                                    <td>{ image.commitHash }</td>
                                    <td>{ image.tag }</td>
                                    <td>{ moment(image.lastBuild).fromNow() }</td>
                                </tr>
                            ) }
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
                
                <Route path="/images/:imageName" component={ ImageDetailModal } />
            </Screen>
        )
    }
}

const mapStateToProps = (state) => ({
    images: state.images.data,
    currentPage: state.images.currentPage,
    totalPages: state.images.totalPages,
    totalImages: state.images.totalImages
});

const mapDispatchToProps = (dispatch) => ({
    onShowPage: (pageNumber) => dispatch(ImagesActions.showPage(pageNumber))
})

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(ImagesScreen));
