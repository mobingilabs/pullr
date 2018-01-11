import * as React from 'react';
import * as moment from 'moment';
import { Icon } from 'react-fa';
import { observer, inject } from 'mobx-react';
import { Route, Link, withRouter, RouteComponentProps } from 'react-router-dom';


import Screen from '../layout/Screen';
import Header from '../layout/Header';
import Button from '../layout/Button';
import Icons from '../layout/Icons';
import Pagination from '../widgets/Pagination';
import Notification from '../widgets/Notification';
import Syncing from '../widgets/Syncing';
import ImageDetailModal from './ImageDetailModal';
import RootStore from '../../state/RootStore';

interface RouteParams {
    imageName?: string;
}

interface Props extends RouteComponentProps<RouteParams> {
    store?: RootStore
}


@withRouter
@inject('store')
@observer
export default class ImagesScreen extends React.Component<Props> {
    actions: [any];
    constructor(props: Props) {
        super(props);

        this.actions = [
            { text: 'Add Image', icon: 'plus', handler: this.handleAddImage }
        ];
    }

    handleAddImage = () => {
        this.props.history.push('/images/add');
    }

    handleEditImage = (imageName: string) => {
        this.props.history.push(`/images/${imageName}/edit`);
    }

    render() {
        const { store } = this.props;
        return (
            <Screen>
                <Header title="IMAGES" subTitle={`${store.images.images.length} images found...`} actions={this.actions} />
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
                            {store.images.images.map(image =>
                                <tr key={image.name}>
                                    <td>
                                        <Link className="table-link" to={`/images/${image.name}`}><Icon name={Icons.images} /> {image.name}</Link>
                                    </td>
                                    <td>{image.sourceProvider}</td>
                                    <td>{image.sourceOwner}/{image.sourceRepository}</td>
                                    <td>{image.builds.map(build => build.tag || build.name).join(', ')}</td>
                                    <td width="100"><Button icon={Icons.edit} onClick={this.handleEditImage.bind(null, image.name)} /></td>
                                </tr>
                            )}

                            <tr>
                                <td colSpan={5}>
                                    <Pagination pagination={store.images.pagination}/>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <Route path="/images/:imageName" exact strict component={ImageDetailModal} />
            </Screen>
        )
    }
}
