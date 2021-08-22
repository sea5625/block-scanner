all: clean
	yarn install
	yarn add --force node-sass caniuse-lite browserslist
	yarn run build

clean:
	yarn cache clean
test:
	CI=true yarn test --watchAll=false
