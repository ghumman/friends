module.exports = {
    "env": {
        "browser": true,
	"node": true,
        "es2021": true
    },
    "extends": [
        "eslint:recommended",
        "plugin:react/recommended"
    ],
    "parserOptions": {
        "ecmaFeatures": {
            "jsx": true
        },
        "ecmaVersion": 12,
        "sourceType": "module"
    },
    "plugins": [
        "react"
    ],
    "rules": {
		"react/prop-types": [0],
		"react/jsx-filename-extension": [0],
		"react/react-in-jsx-scope": [0],
		"react/jsx-tag-spacing": [0],
		"react/jsx-space-before-closing": [0],
		"react/prefer-stateless-function": [0],
		"import/extensions": [0],
		"import/no-unresolved": [0],
		"space-before-function-paren": [0],
		"class-methods-use-this": [0],
		"import/no-named-as-default": [0],
		"import/no-extraneous-dependencies": [0],
		"global-require": [0],
		"no-param-reassign": [0],
		"react/require-default-props": [0],
		"react/forbid-prop-types": [0],
		"react/no-unused-prop-types": [1],
		"jsx-a11y/href-no-hash": "off",
		"eqeqeq": 2
	}
};
