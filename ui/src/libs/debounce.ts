type Func<R, A1, A2, A3, A4, A5, A6, A7, A8, A9> = (a1?: A1, a2?: A2, a3?: A3, a4?: A4, a5?: A5, a6?: A6, a7?: A7, a8?: A8, a9?: A9) => R;
export default function debounce<R, A1 = never, A2 = never, A3 = never, A4 = never, A5 = never, A6 = never, A7 = never, A8 = never, A9 = never>
    (delay: number, func: Func<R, A1, A2, A3, A4, A5, A6, A7, A8, A9>): Func<void, A1, A2, A3, A4, A5, A6, A7, A8, A9> {
    let timeoutId: number = null;

    return (a1?: A1, a2?: A2, a3?: A3, a4?: A4, a5?: A5, a6?: A6, a7?: A7, a8?: A8, a9?: A9) => {
        if (timeoutId !== null) {
            clearTimeout(timeoutId);
        }

        timeoutId = setTimeout(() => {
            func(a1, a2, a3, a4, a5, a6, a7, a8, a9);
        }, delay);
    };
}
