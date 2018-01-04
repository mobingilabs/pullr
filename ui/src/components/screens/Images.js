import React from 'react';
import { Route, Link, withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { Icon } from 'react-fa';
import moment from 'moment';

import Screen from '../layout/Screen';
import Header from '../layout/Header';
import Button from '../layout/Button';
import Icons from '../layout/Icons';
import Pagination from '../widgets/Pagination';
import Notification from '../widgets/Notification';
import ImageDetailModal from './ImageDetailModal';

import * as ImagesActions from '../../state/images/actions';

export class ImagesScreen extends React.PureComponent {
    constructor(props) {
        super(props);

        this.actions = [
            { text: 'Add Image', icon: 'plus', handler: this.handleAddImage }
        ];
    }

    handleAddImage = () => {
        this.props.history.push('/images/add');
    }

    handleEditImage = (imageName) => {
        this.props.history.push(`/images/${imageName}/edit`);
    }

    render() {
        const { images, imageOrder } = this.props;
        return (
            <Screen>
                <Header title="IMAGES" subTitle={`${this.props.totalImages} images found...`} actions={this.actions} />
                <Notification id="images-create" />
                <div className="scroll-horizontal">
                    <table className="wide big-shadow">
                        <thead>
                            <tr>
                                <th>IMAGE NAME</th>
                                <th>SOURCE PROVIDER</th>
                                <th>REPOSITORY</th>
                                <th>TAGS</th>
                                <th></th>
                            </tr>
                        </thead>
                        <tbody>
                            {imageOrder.map(key =>
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

                            <tr>
                                <td colSpan="5">
                                    <Pagination
                                        currentPage={this.props.currentPage}
                                        totalPages={this.props.totalPages}
                                        onChangePage={this.props.onShowPage}
                                        itemPerPage={10}
                                        totalItems={this.props.totalImages} />
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <Route path="/images/:imageName" component={ImageDetailModal} />
            </Screen>
        )
    }
}

const mapStateToProps = (state) => ({
    images: state.images.data,
    imageOrder: state.images.dataOrder,
    currentPage: state.images.currentPage,
    totalPages: state.images.totalPages,
    totalImages: state.images.totalImages
});

const mapDispatchToProps = (dispatch) => ({
    onShowPage: (pageNumber) => dispatch(ImagesActions.showPage(pageNumber))
})

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(ImagesScreen));
