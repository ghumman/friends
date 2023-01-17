## Following commands were used to create this project
```
npx create-react-app react-js
cd react-js
npm start
```

## Use following command to push the code
```
If on branch react-js: 
git push react-js react-js
If incase on branch master: 
git push react-js master:react-js
```

## Initial push done using following commands
```
git remote add react-js git@github.com:ghumman/friends.git
git push react-js master:react-js
```

## Eslint included using following commands
```
yarn add eslint --dev
yarn run eslint --init
yarn run eslint yourfile.js
```

## When Eslint is giving error 'error: Environment key "es2021" is unknown'
```
yarn add eslint --save-dev
./node_modules/.bin/eslint --init
yarn run esline . or ./node_modules/.bin/eslint yourfile.js

```

## When Eslint starts running, we start getting the following error doing 'yarn start'
```
yarn run v1.22.5
$ react-scripts start

There might be a problem with the project dependency tree.
It is likely not a bug in Create React App, but something you need to fix locally.

The react-scripts package provided by Create React App requires a dependency:

  "eslint": "^6.6.0"

Don't try to install it manually: your package manager does it automatically.
However, a different version of eslint was detected higher up in the tree:

  /home/ghumman/programming/friends/react-js/node_modules/eslint (version: 7.11.0) 

Manually installing incompatible versions is known to cause hard-to-debug issues.

If you would prefer to ignore this check, add SKIP_PREFLIGHT_CHECK=true to an .env file in your project.
That will permanently disable this message but you might encounter other issues.

To fix the dependency tree, try following the steps below in the exact order:

  1. Delete package-lock.json (not package.json!) and/or yarn.lock in your project folder.
  2. Delete node_modules in your project folder.
  3. Remove "eslint" from dependencies and/or devDependencies in the package.json file in your project folder.
  4. Run npm install or yarn, depending on the package manager you use.

In most cases, this should be enough to fix the problem.
If this has not helped, there are a few other things you can try:

  5. If you used npm, install yarn (http://yarnpkg.com/) and repeat the above steps with it instead.
     This may help because npm has known issues with package hoisting which may get resolved in future versions.

  6. Check if /home/ghumman/programming/friends/react-js/node_modules/eslint is outside your project directory.
     For example, you might have accidentally installed something in your home folder.

  7. Try running npm ls eslint in your project folder.
     This will tell you which other package (apart from the expected react-scripts) installed eslint.

If nothing else helps, add SKIP_PREFLIGHT_CHECK=true to an .env file in your project.
That would permanently disable this preflight check in case you want to proceed anyway.

P.S. We know this message is long but please read the steps above :-) We hope you find them helpful!

error Command failed with exit code 1.
info Visit https://yarnpkg.com/en/docs/cli/run for documentation about this command.
```

## Fixing eslinting issue
Following above instructions created a file '.env' and added following line. 
```
SKIP_PREFLIGHT_CHECK=true
```
Able to run `yarn start` after that. Also able to `yarn run eslint .`

## Final Comments about eslint 
After all this ESLint plugin installed in VSCode start throwing errors. So I uninstalled it reloaded and install the plugin again.

## ESLinting rules
Copied eslinting rules from another project to .eslintrc.js which fixed all props validation errors. The same file is used by both command line `yarn run eslint .` and vscode eslint plugin.

## Deploying React App to Github Pages
This app is hosted on github pages and can be accessed at `https://ghumman.github.io/friends/`. Used following two links to figure out how to deploy on github pages. 
```
https://github.com/gitname/react-gh-pages
https://docs.github.com/en/pages/getting-started-with-github-pages/about-github-pages
```
Here's summary.
 
There are 2 main steps.

- Step 1 : Coding

Install npm package. 
```
npm install gh-pages --save-dev
```
In `package.json` after version added following line. 
```
"homepage": "https://ghumman.github.io/friends",
``` 
And before line `start: react-scripts start`, add following lines
```
"predeploy": "npm run build",
"deploy": "gh-pages -d build",
``` 
Then run following command.
```
npm run deploy
``` 
Acccording to the web link provided above, this will run predeploy and deploy commands, create a new branch gh-pages in your repository and push your build directory to this newly created `gh-pages` branch. If you already have that branch it will update it. You will have a new commit with default message `Updates`. If you want to have a different message use `npm run deploy -- -m "Deploy React app to GitHub Pages"`.

Note: When I ran it, it did not work and it was complaining about url not set. I had to run following command to fix it. This is the case when you have ssh communication set up. 
```
git remote set-url origin git@github.com:ghumman/friends.git
```

- Step 2: Github Settings

Go to your github repo settings and then to `Pages`. Select `Source` as `Deploy from a branch`, `Branch` as `gh-pages`, `Folder` as `/root`. Then save it. You can check the deployment status inside `Actions`.

---------------------------------------------------------------------------------
Automatically Generated Documentation
---------------------------------------------------------------------------------


This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

## Available Scripts

In the project directory, you can run:

### `yarn start`

Runs the app in the development mode.<br />
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.<br />
You will also see any lint errors in the console.

### `yarn test`

Launches the test runner in the interactive watch mode.<br />
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

### `yarn build`

Builds the app for production to the `build` folder.<br />
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.<br />
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

### `yarn eject`

**Note: this is a one-way operation. Once you `eject`, you can’t go back!**

If you aren’t satisfied with the build tool and configuration choices, you can `eject` at any time. This command will remove the single build dependency from your project.

Instead, it will copy all the configuration files and the transitive dependencies (webpack, Babel, ESLint, etc) right into your project so you have full control over them. All of the commands except `eject` will still work, but they will point to the copied scripts so you can tweak them. At this point you’re on your own.

You don’t have to ever use `eject`. The curated feature set is suitable for small and middle deployments, and you shouldn’t feel obligated to use this feature. However we understand that this tool wouldn’t be useful if you couldn’t customize it when you are ready for it.

## Learn More

You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).

To learn React, check out the [React documentation](https://reactjs.org/).

### Code Splitting

This section has moved here: https://facebook.github.io/create-react-app/docs/code-splitting

### Analyzing the Bundle Size

This section has moved here: https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size

### Making a Progressive Web App

This section has moved here: https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app

### Advanced Configuration

This section has moved here: https://facebook.github.io/create-react-app/docs/advanced-configuration

### Deployment

This section has moved here: https://facebook.github.io/create-react-app/docs/deployment

### `yarn build` fails to minify

This section has moved here: https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify
