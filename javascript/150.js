/**
 * @param {string[]} tokens
 * @return {number}
 */
var evalRPN = function(tokens) {
    const set = new Set(["+", "-", "*", "/"]);

    const stack = [];

    for (let ch of tokens) {
        if (set.has(ch)) {
            let result  = 0;
            let second = stack.pop();
            let first = stack.pop();

            switch (ch) {
                case '+':
                    result = first + second; 
                    break;
                case '-':
                    result = first - second; 
                    break;
                case '*':
                    result = first * second; 
                    break;
                case '/':
                    result = Math.trunc(first / second); 
                    break;
            }
            stack.push(result);
        } else {
            stack.push(Number(ch));
        }
    }

    return stack[stack.length - 1];
};