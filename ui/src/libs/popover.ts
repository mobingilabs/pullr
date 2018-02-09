declare module 'react-popover' {
    import * as React from 'react';
    type Place = 'above' | 'right' | 'below' | 'left' | 'row' | 'column' | 'start' | 'end';
    interface Props {
        className?: string;
        body: React.ReactNode | Array<React.ReactNode>;
        isOpen: boolean;
        preferPlace?: Place;
        place?: Place;
        onOuterAction: (e: Event) => void;
        appendTarget: HTMLElement;
    }

    class Popover extends React.Component<Props> { }
    export = Popover;
}
