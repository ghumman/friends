class Solution {
    public int[] productExceptSelf(int[] nums) {
        /*
        [1,2,3,4]
        L= [1, 1, 2, 6]
        R= [24, 12, 4, 1]
        */

        int[] product = new int[nums.length];
        product[0] = 1; 
        for ( int i =1; i<nums.length; i++) {
            product[i] = product[i-1]* nums[i-1];
        }
        int suffix = 1; 
        for (int i= nums.length - 1; i>= 0; i--) {
            product[i] *= suffix; 
            suffix *= nums[i];
        }

        return product; 
    } 
}