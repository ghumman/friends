```
git init
git remote add angular git@github.com:ghumman/friends.git
git add .
git commit -m "Initial commit"
git push angular master:angular
```

## Run Project
`ng serve` or `ng serve --open` to run the app and open the browser 

## Create Component
```
ng generate component Login
```

## Deploying This App On Github Pages
Although this application is part of friends but it is deployed on repo `https://github.com/ghumman/friends-angular-ui`. This is because only one github pages application can be lauched from one repository. React application is deployed from `https://github.com/ghumman/friends` and Vue application is deployed from `https://github.com/ghumman/friends-vue-js`. For general github web application deployed, do checkout React branch for more documentation. Following are the steps I did to deploy an angualar application to github pages. 

Run following
```
git remote add origin git@github.com:ghumman/friends-angular-ui.git
ng add angular-cli-ghpages
ng deploy --base-href=/friends-angular-ui/
``` 

Somehow when you deploy, it knows that this should take remote address `origin` and push it to branch `gh-pages` inside `origin` repo.

Now this repo has two remote addresses. `origin` which points to `friends-angular-ui` and `angular` which points to `friends` repo. Just make sure before pushing where exactly you're pushing things. 

//////////////////////////////////////////////////
//  Auto Generated
//////////////////////////////////////////////////
# Angular

This project was generated with [Angular CLI](https://github.com/angular/angular-cli) version 11.0.3.

## Development server

Run `ng serve` for a dev server. Navigate to `http://localhost:4200/`. The app will automatically reload if you change any of the source files.

## Code scaffolding

Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

## Build

Run `ng build` to build the project. The build artifacts will be stored in the `dist/` directory. Use the `--prod` flag for a production build.

## Running unit tests

Run `ng test` to execute the unit tests via [Karma](https://karma-runner.github.io).

## Running end-to-end tests

Run `ng e2e` to execute the end-to-end tests via [Protractor](http://www.protractortest.org/).

## Further help

To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI Overview and Command Reference](https://angular.io/cli) page.
