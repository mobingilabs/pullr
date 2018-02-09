import * as React from 'react';
import * as moment from 'moment';
import { Icon } from 'react-fa';
import { computed } from 'mobx';
import { observer, inject } from 'mobx-react';
import { Route, Link, withRouter, RouteComponentProps } from 'react-router-dom';
import * as fuzzysearch from 'fuzzysearch';

import Image from '../../state/models/Image';
import ImagesStore from '../../state/ImagesStore';
import TableOptions from '../../state/models/TableOptions';

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
import InProgress from '../widgets/InProgress';
import LoadingSpinner from '../widgets/LoadingSpinner';

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

    options: TableOptions;

    constructor(props: Props) {
        super(props);

        this.options = new TableOptions();
        this.popoverPortal = document.getElementById('popover');
        this.actions = [
            { text: 'Add Image', icon: 'plus', handler: this.handleAddImage }
        ];
    }

    @computed
    get images(): Image[] {
        if (this.options.query.length > 0) {
            return this.props.images.images.filter(img => fuzzysearch(this.options.query, img.name))
        }

        return this.props.images.images;
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
                                    <th className="table-actions-column">
                                        <TableActionsMenu
                                            tableOptions={this.options}
                                            listOptions={this.props.images.listOptions}
                                            sortableFields={{ name: 'Name', created_at: 'Created At' }}
                                            portal={this.popoverPortal} />
                                    </th>
                                </tr>
                            </thead>
                            <tbody>
                                <InProgress cmd={this.props.images.fetchImages}>
                                    <tr>
                                        <td colSpan={5} className="table-loading-row">
                                            <LoadingSpinner /> Loading images...
                                        </td>
                                    </tr>
                                </InProgress>
                                {this.images.map(image =>
                                    <tr key={image.name}>
                                        <td>
                                            <Link className="table-link" to={`/images/${image.key}`}><Icon name={Icons.images} /> {image.name}</Link>
                                        </td>
                                        <td>{image.repository.provider}</td>
                                        <td>{image.repository.owner}/{image.repository.name}</td>
                                        <td>{image.tags.map(tag => tag.ref_test || tag.name).join(', ')}</td>
                                        <td width={100}><Button icon={Icons.edit} onClick={this.handleEditImage.bind(null, image.key)} /></td>
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
                <Route path="/images/:imageKey" exact strict component={ImageDetailModal} />
            </Screen>
        )
    }
}
