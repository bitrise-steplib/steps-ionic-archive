# Ionic archive

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/steps-ionic-archive?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/steps-ionic-archive/releases)

Generates builds for the requested platforms.

<details>
<summary>Description</summary>

Generates an iOS, Android or an iOS and Android build based on your Step settings and a `build.json` file which is inherited from the **Generate cordova build configuration** Step.

### Configuring the Step
1. In the **Platform to use in ionic-cli commands** input, select the platform you wish to build for from the drop-down menu.
2. In the **Build command configuration** input, select either `debug` or `release` mode from the drop-down menu. If you set `release` to generate a release build, then you should select `device` instead of `emulator` in the **Build command target**.
3. Select either `emulator` or `device` in the **Build command target** field.
4. Make sure the file path set in the **Working directory** is pointing to the path where your project's code got downloaded to.
If you're using any services of the Ionic framework which requires login credentials, you have to provide those in the **Ionic username** and **Ionic password** inputs. Don't worry, these are stored as secret Environment Variables.
If you wish to to modify your native projects by inserting the **Ionic Prepare** Step before the **Ionic Archive** Step, then make sure you set the **Should `ionic cordova prepare` be executed before `ionic cordova build`?** input to `false` in the **Ionic Archive** Step.
The **Build configuration path to describe code signing properties** input is automatically filled out with the output Environment Variable of the **Generate cordova build configuation** Step. You do not need to modify this input.

### Troubleshooting
Make sure you add the **Ionic Archive** Step AFTER the **Generate cordova build configuration** Step as the latter generates the build config file which the **Ionic Archive** Step uses to successfully build an iOS and/or Android project.
Make sure you insert the **Ionic Archive** Step BEFORE any deploy Step.
Note that if you’re building for both iOS and Android in one workflow, and either of your apps fails, the whole **Ionic Archive** Step will fail.
If you set the **Build configuration** input in the **Generate Cordova Build Configuration** Step to `release`, then you need to use the release configuration in the **Ionic Archive** Step as well.

### Useful links
- [Getting started with Ionic/Cordova apps on Bitrise](https://devcenter.bitrise.io/code-signing/android-code-signing/android-code-signing-using-bitrise-sign-apk-step/)
- [Secret Environment Variables on Bitrise](https://devcenter.bitrise.io/builds/env-vars-secret-env-vars/)
- [Ionic framework](https://ionicframework.com/)

### Related Steps
- [Generate Cordova Build Configuration](https://www.bitrise.io/integrations/steps/android-build)
- [Run npm command](https://www.bitrise.io/integrations/steps/npm)
- [Cordova Archive](https://www.bitrise.io/integrations/steps/cordova-archive)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `platform` | Specify this input to apply ionic-cli commands to desired platforms only.  `ionic cordova build [OTHER_PARAMS] <platform>` | required | `ios,android` |
| `configuration` | Specify build command configuration.  `ionic cordova build [OTHER_PARAMS] [--release \| --debug]` | required | `release` |
| `target` | Specify build command target.  `ionic cordova build [OTHER_PARAMS] [--device \| --emulator]` | required | `device` |
| `build_config` | Path to the build configuration file (build.json), which describes code signing properties. |  | `$BITRISE_CORDOVA_BUILD_CONFIGURATION` |
| `options` | Use this input to specify custom options, to append to the end of the ionic-cli build command.  Cordova now supports the new build system made default in XCode 10 (https://github.com/apache/cordova-ios/issues/407). To use the legacy build system add `-- --buildFlag="-UseModernBuildSystem=0"` to the options string.  Example: - `--browserify`  `ionic cordova build [OTHER_PARAMS] [options]` |  |  |
| `ionic_username` | Use `Ionic username` and `Ionic password` to login with ionic-cli. | sensitive |  |
| `ionic_password` | Use `Ionic username` and `Ionic password` to login with ionic-cli. | sensitive |  |
| `ionic_version` | The version of ionic you want to use.  If value is set to `latest`, the step will update to the latest ionic version. Leave this input empty to use the preinstalled ionic version. |  |  |
| `run_ionic_prepare` | It should be set to false if ionic-prepare step is used.  - false: `ionic cordova build` - true: `ionic cordova prepare --no-build` followed by `ionic cordova build` |  | `true` |
| `cordova_version` | The version of cordova you want to use.  If value is set to `latest`, the step will update to the latest cordova version. Leave this input empty to use the preinstalled cordova version. |  |  |
| `workdir` | Root directory of your Ionic project, where your Ionic config.xml exists. | required | `$BITRISE_SOURCE_DIR` |
| `android_app_type` | Set the distribution type that you want to build for your Android app.  | required | `apk` |
| `cache_local_deps` | Select if the contents of node_modules directory should be cached. `true`: Mark local dependencies to be cached. `false`: Do not use cache.  | required | `false` |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `BITRISE_IPA_PATH` |  |
| `BITRISE_APP_DIR_PATH` |  |
| `BITRISE_APP_PATH` |  |
| `BITRISE_DSYM_DIR_PATH` |  |
| `BITRISE_DSYM_PATH` |  |
| `BITRISE_APK_PATH` |  |
| `BITRISE_APK_PATH_LIST` |  |
| `BITRISE_AAB_PATH` | This output will include the path of the generated AAB. If the build generates more than one AAB this output will contain the last one's path. |
| `BITRISE_AAB_PATH_LIST` | This output will include the paths of the generated AABs. The paths are separated with `\|` character, for example, `app--debug.aab\|app-mips-debug.aab` |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/steps-ionic-archive/pulls) and [issues](https://github.com/bitrise-steplib/steps-ionic-archive/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

**Note:** this step's end-to-end tests (defined in `e2e/bitrise.yml`) are working with secrets which are intentionally not stored in this repo. External contributors won't be able to run those tests. Don't worry, if you open a PR with your contribution, we will help with running tests and make sure that they pass.

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
