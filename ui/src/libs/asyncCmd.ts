import * as Promise from 'bluebird';
import { observable, action, transaction } from 'mobx';

import ApiError from './api/ApiError';

type Func<R, A1, A2, A3, A4, A5, A6, A7, A8, A9> = (a1?: A1, a2?: A2, a3?: A3, a4?: A4, a5?: A5, a6?: A6, a7?: A7, a8?: A8, a9?: A9) => R;
export default class AsyncCmd<T, E = ApiError, A1 = never, A2 = never, A3 = never, A4 = never, A5 = never, A6 = never, A7 = never, A8 = never, A9 = never> {
    @observable inProgress = false;
    @observable err: E;
    @observable value: T;

    run: Func<Promise<void>, A1, A2, A3, A4, A5, A6, A7, A8, A9>;

    private handler: Func<Promise<T>, A1, A2, A3, A4, A5, A6, A7, A8, A9>;
    private canceler: Function;

    constructor(handler: Func<Promise<T>, A1, A2, A3, A4, A5, A6, A7, A8, A9>) {
        this.handler = handler;

        this.run = (a1, a2, a3, a4, a5, a6, a7, a8, a9): Promise<void> => {
            let canceled = false;
            transaction(() => {
                if (this.canceler) {
                    this.canceler();
                }

                this.startProsessing();
            });

            this.canceler = () => {
                canceled = true;
                this.stopProcessing();
            };

            return this.handler(a1, a2, a3, a4, a5, a6, a7, a8, a9)
                .then((val: T) => {
                    if (!canceled) this.handleSuccess(val);
                })
                .catch((err: E) => !canceled && this.handleFailure(err))
                .finally(() => !canceled && this.handleFinish());
        }
    }

    @action.bound
    private startProsessing() {
        this.inProgress = true;
    }

    @action.bound
    private stopProcessing() {
        this.inProgress = false;
    }

    @action.bound
    handleSuccess(val: T) {
        this.value = val;
    }

    @action.bound
    handleFailure(err: E) {
        this.err = err;
    }

    @action.bound
    handleFinish() {
        this.inProgress = false;
    }
}
