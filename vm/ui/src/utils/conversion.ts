export enum Units {
    Gigabyte = 1,
    Bytes = 2,
}

export const convert = (value: number, from: Units, to: Units): number => {
    if (from == Units.Gigabyte && to == Units.Bytes) {
        return value * 1e9; // 1 GB = 10^9 Bytes
    }

    if (from == Units.Bytes && to == Units.Gigabyte) {
        return value / 1e9; // 1 GB = 10^9 Bytes
    }

    return 0;
};
