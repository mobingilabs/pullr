import { observable } from 'mobx';
import ListOptions from './ListOptions';

export default class TableOptions {
    @observable query: string = '';
}
