'use strict';

// Copyright (c) 2017 hirowaki https://github.com/hirowaki

// making static HTML documents.

// modules.
// https://www.npmjs.com/package/fs-extra
// https://www.npmjs.com/package/fs-extra-promise
// https://www.npmjs.com/package/metalsmith
const Promise = require("bluebird");
const ejs = require('ejs');
const path = require('path');
const fs = require('fs-extra-promise').usePromise(Promise); // using bluebird.
const Metalsmith = require('metalsmith');
const markdown = require('metalsmith-markdown');
const layouts = require('metalsmith-layouts');
const exec = require('child_process').exec;

// Constants.
const scriptFolder = './_scripts';
const srcFolder = './_src';
const dstFolder = './_build';
const libFolder = './_lib';
const cssFolder = './_css';
const assetFolder = './assets';

const MD_INDEX = `---
layout: index.ejs
directories: [DIRECTORIES]
---`;

const MD_PAGE = `---
layout: sample.ejs
subtitle: SUBTITLE
---
SCRIPT
`;

// have to add <script> tag. Otherwise, 'SCRIPT' section will be formatted as HTML. :(
const MD_SOURCE = `---
layout: code.ejs
subtitle: SUBTITLE
---
<script>
SCRIPT
</script>
`;


// preProcess.
function preProcess() {
    // clean up work folders just in case.
    return Promise.all([srcFolder, scriptFolder].map((folder) => {
        return fs.removeAsync(folder);
    }));
}

// enumerate go sample directories.
function enumGoFiles() {
    return fs.readdirAsync('./src')
    .then((files) => {
        const directories = files.filter((filepath) => {
            filepath = './src/' + filepath;

            if (fs.statSync(filepath).isDirectory()) {
                return /sample_.*$/.test(filepath);
            }
            return false;
        });
        directories.sort((a, b) => {
            return a.name - b.name;
        });
        return directories;
    });
}

// enumerate go sample directories.
function gopherJS(sampleFolders) {
    function _execGopherJS(dirname) {
        // 1. copy main.go to _scripts.
        // 2. run GopherJS.
        // 3, rename main.js to sample_*.js.

        const srcGo = `./${scriptFolder}/main.go`;
        const desJS = `./${scriptFolder}/${dirname}.js`;

        return fs.copyAsync(`./src/${dirname}/main.go`, srcGo)
        .then(() => {
            return new Promise((resolve, reject) => {
                exec(`gopherjs build ${srcGo} -o ${desJS}`, (err, stdout, stderr) => {
                    void(stdout);
                    void(stderr);

                    if (err) {
                        return reject(err);
                    }
                    return resolve();
                });
            });
        })
        .then(() => {
            return fs.removeAsync(srcGo);
        });
    }

    return sampleFolders.reduce((p, dirname) => {
        return p.then(() => {
            return _execGopherJS(dirname);
        });
    }, Promise.resolve())
}

function makeSourceFiles(samples) {
    function makeIndex() {
        let data = MD_INDEX.replace('DIRECTORIES', samples.toString());

        return fs.outputFileAsync(`${srcFolder}/index.md`, data);
    }

    function makeEachPage() {
        return Promise.resolve(samples)
        .each((sample) => {
            // make each one.
            const targetDir = `${srcFolder}/${sample}`;
            const go = `./src/${sample}/main.go`;

            return Promise.all([
                fs.ensureDirAsync(targetDir),
                fs.readFileAsync(go)
            ])
            .then((res) => {
                let demo = MD_PAGE.replace('SUBTITLE', sample)
                            .replace('SCRIPT', `<script src="./${sample}.js"></script>`);
                let code = MD_SOURCE.replace('SUBTITLE', sample).replace('SCRIPT', res[1].toString());

                const js =  `${sample}.js`;

                return Promise.all([
                    fs.outputFileAsync(`${targetDir}/index.md`, demo),
                    fs.outputFileAsync(`${targetDir}/code.md`, code),
                    fs.copyAsync(`${scriptFolder}/${js}`, `${targetDir}/${js}`),
                    fs.copyAsync(`${scriptFolder}/${js}.map`, `${targetDir}/${js}.map`)
                ]);
            });
        });
    }

    return Promise.all([
        makeIndex(),
        makeEachPage()
    ]);
}

// Metalsmith.
// https://www.npmjs.com/package/metalsmith
function runMetalsmith() {
    return new Promise((resolve, reject) => {
        Metalsmith(__dirname)
        .metadata({
            title: "ebiten tutorials"
        })
        .source(srcFolder)          // set source folder.
        .destination(dstFolder)     // set destination folder.
        .clean(true)                // clean up the destination folder.
        .use(markdown())            // markdown.
        .use(layouts({              // set up layout engine.
            engine: 'ejs',
            directory: '_layouts'
        }))
        .build(function (err, files) {  // run.
            void(files);
            if (err) {
                return reject(err);
            }
            return resolve();
        });
    });
}

function postProcess() {
    return Promise.all([
        fs.copyAsync(`${libFolder}`, `${dstFolder}/lib`),
        fs.copyAsync(`${assetFolder}`, `${dstFolder}/${assetFolder}`),
        fs.copyAsync(`${cssFolder}`, `${dstFolder}/css`)
    ])
    .then(() => {
        // clean up work folders.
        return Promise.all([srcFolder, scriptFolder].map((folder) => {
            return fs.removeAsync(folder);
        }));
    });
}

// main.
function main() {
    return preProcess()
    .then(() => {
        return enumGoFiles();
    })
    .then((samples) => {
        return gopherJS(samples)
        .then(() => {
            return makeSourceFiles(samples);
        })
        .then(() => {
            return runMetalsmith();
        })
        .then(() => {
            return postProcess();
        });
    })
    .catch((err) => {
        console.log(err);
    });
}

main();
