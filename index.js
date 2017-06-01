'use strict';

// Copyright (c) 2017 hirowaki https://github.com/hirowaki

// making static HTML documents.
// I used to use MetalSmith to build HTML, bit I figured out just using ejs was enough.

const Promise = require("bluebird");
const ejs = Promise.promisifyAll(require('ejs'));
const path = require('path');
const fs = require('fs-extra-promise').usePromise(Promise); // using bluebird.
const exec = require('child_process').exec;
const EnumFiles = require('enum-files');

// Constants.
const scriptFolder = './_scripts';
const dstFolder = './_build';
const libFolder = './_lib';
const cssFolder = './_css';
const assetFolder = './assets';

// preProcess.
function preProcess() {
    // clean up work folders just in case.
    return fs.removeAsync(scriptFolder);
}

// enumerate go sample directories.
function enumGoFiles() {
    return EnumFiles.dir('./src')
    .then((directories) => {
        return directories.filter((directory) => {
            return /sample_.*$/.test(directory);
        })
        .map(directory => path.basename(directory));
    });
}

// enumerate go sample directories.
function gopherJS(sampleFolders) {
    function _execGopherJS(dirname) {
        // 1. copy main.go to _scripts.
        // 2. run GopherJS.
        // 3, rename main.js to sample_*.js.

        const srcGo = `${scriptFolder}/main.go`;
        const dstJS = `${scriptFolder}/${dirname}.js`;

        return fs.copyAsync(`./src/${dirname}/main.go`, srcGo)
        .then(() => {
            return new Promise((resolve, reject) => {
                exec(`gopherjs build ${srcGo} -o ${dstJS}`, (err, stdout, stderr) => {
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

function processEjs(dst, src, data) {
    return ejs.renderFileAsync(src, data)
    .then((html) => {
        return fs.outputFileAsync(dst, html);
    });
}

function makeIndex(samples) {
    const dst = `${dstFolder}/index.html`;
    const src = '_layouts/index.ejs';

    return processEjs(dst, src, {
        title: "ebiten samples",
        directories: samples
    });
}

function makeSourceFiles(samples) {
    return Promise.resolve(samples)
    .each((sample) => {
        // make each one.
        const targetDir = `${dstFolder}/${sample}`;
        const go = `./src/${sample}/main.go`;

        return Promise.all([
            fs.ensureDirAsync(targetDir),
            fs.readFileAsync(go)
        ])
        .then((res) => {
            const js =  `${sample}.js`;

            const data = {
                title: "ebiten samples",
                subtitle: sample,
                js: `./${js}`,
                code: res[1].toString()
            };

            return Promise.all([
                processEjs(`${targetDir}/index.html`, '_layouts/sample.ejs', data),
                processEjs(`${targetDir}/code.html`, '_layouts/code.ejs', data),
                fs.copyAsync(`${scriptFolder}/${js}`, `${targetDir}/${js}`),
                fs.copyAsync(`${scriptFolder}/${js}.map`, `${targetDir}/${js}.map`)
            ]);
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
        return fs.removeAsync(scriptFolder);
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
            return Promise.all([
                makeIndex(samples),
                makeSourceFiles(samples),
            ]);
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
