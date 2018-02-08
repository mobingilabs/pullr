import * as React from 'react';
import * as moment from 'moment';
import { Icon } from 'react-fa';
import { observer, inject } from 'mobx-react';
import { Route, Link, withRouter, RouteComponentProps } from 'react-router-dom';

import ImagesStore from '../../state/ImagesStore';

import Screen from '../layout/Screen';
import Header from '../layout/Header';
import Button from '../layout/Button';
import Icons from '../layout/Icons';
import Pagination from '../widgets/Pagination';
import Notification from '../widgets/Notification';
import Syncing from '../widgets/Syncing';
import ImageDetailModal from './ImageDetailModal';
import ApiError from '../../libs/api/ApiError';
import Alert from '../widgets/Alert';
import TableActionsMenu from '../widgets/TableActionsMenu';

interface RouteParams {
    imageName?: string;
}

interface Props extends RouteComponentProps<RouteParams> {
    images?: ImagesStore;
}

interface State {
    showLoadErr: boolean;
}

@withRouter
@inject('images')
@observer
export default class ImagesScreen extends React.Component<Props, State> {
    actions: [any];
    popoverPortal: HTMLElement;
    constructor(props: Props) {
        super(props);

        this.popoverPortal = document.getElementById('popover');
        this.actions = [
            { text: 'Add Image', icon: 'plus', handler: this.handleAddImage }
        ];
    }

    componentDidMount() {
        this.props.images.fetchImages.run().done();
    }

    handleAddImage = () => {
        this.props.history.push('/images/add');
    }

    handleEditImage = (imageName: string) => {
        this.props.history.push(`/images/${imageName}/edit`);
    }

    handleGotoPage = (page: number) => {
        this.props.images.fetchImagesAtPage(page);
    }

    render() {
        const showLoadErr = !!this.props.images.fetchImages.err;
        return (
            <Screen>
                <Header title="IMAGES" subTitle={!showLoadErr && `${this.props.images.images.length} images found...`} actions={this.actions} />
                <Notification id="images-create" />
                {showLoadErr && <Alert message="An error occured while loading images, please try again later." />}
                {!showLoadErr &&
                    <div className="scroll-horizontal">
                        <table className="wide big-shadow">
                            <thead>
                                <tr>
                                    <th>IMAGE NAME</th>
                                    <th>SOURCE PROVIDER</th>
                                    <th>REPOSITORY</th>
                                    <th>TAGS</th>
                                    <th><TableActionsMenu listOptions={this.props.images.listOptions} portal={this.popoverPortal} /></th>
                                </tr>
                            </thead>
                            <tbody>
                                {this.props.images.images.map(image =>
                                    <tr key={image.name}>
                                        <td>
                                            <Link className="table-link" to={`/images/${image.name}`}><Icon name={Icons.images} /> {image.name}</Link>
                                        </td>
                                        <td>{image.repository.provider}</td>
                                        <td>{image.repository.owner}/{image.repository.name}</td>
                                        <td>{image.tags.map(tag => tag.ref_test || tag.name).join(', ')}</td>
                                        <td width={100}><Button icon={Icons.edit} onClick={this.handleEditImage.bind(null, image.name)} /></td>
                                    </tr>
                                )}
                                {this.props.images.pagination.last != 0 &&
                                    <tr>
                                        <td colSpan={5}>
                                            <Pagination onGotoPage={this.handleGotoPage} pagination={this.props.images.pagination} />
                                        </td>
                                    </tr>
                                }
                            </tbody>
                        </table>
                    </div>
                }
                <Route path="/images/:imageName" exact strict component={ImageDetailModal} />
            </Screen>
        )
    }
}
