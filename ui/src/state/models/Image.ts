import { observable, computed } from "mobx";
import { IImageBuild, ImageBuild } from "./ImageBuild";

export interface IImage {
    name: string;
    sourceProvider: string;
    sourceOwner: string;
    sourceRepository: string;
    dockerfilePath: string;
    builds: Array<IImageBuild>
}

export default class Image implements IImage {
    @observable name: string;
    @observable sourceProvider: string;
    @observable sourceOwner: string;
    @observable sourceRepository: string;
    @observable dockerfilePath: string;
    @observable builds: Array<ImageBuild>;

    constructor(json: IImage) {
        this.name = json.name;
        this.sourceProvider = json.sourceProvider;
        this.sourceOwner = json.sourceOwner;
        this.sourceRepository = json.sourceRepository;
        this.dockerfilePath = json.dockerfilePath;
        this.builds = json.builds.map(build => new ImageBuild(build));
    }

    addBuild() {
        this.builds.push(ImageBuild.create());
    }

    removeBuild(buildIndex: number) {
        this.builds.splice(buildIndex, 1);
    }

    clone(): Image {
        return new Image({
            name: this.name,
            sourceProvider: this.sourceProvider,
            sourceOwner: this.sourceOwner,
            sourceRepository: this.sourceRepository,
            dockerfilePath: this.dockerfilePath,
            builds: this.builds.map(build => build.clone())
        });
    }

    static create(): Image {
        return new Image({
            name: '',
            sourceProvider: '',
            sourceOwner: '',
            sourceRepository: '',
            dockerfilePath: 'Dockerfile',
            builds: [
                new ImageBuild({ type: 'branch', name: 'master', tag: 'latest' })
            ]
        });
    }
}
