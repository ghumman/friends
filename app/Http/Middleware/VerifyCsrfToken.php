<?php

namespace App\Http\Middleware;

use Illuminate\Foundation\Http\Middleware\VerifyCsrfToken as Middleware;

class VerifyCsrfToken extends Middleware
{
    /**
     * The URIs that should be excluded from CSRF verification.
     *
     * @var array
     */
    protected $except = [
        '/add-user',
        '/login',
        '/change-password',
        '/forgot-password',
        '/reset-password',
        '/all-friends',
        '/send-message',
        '/messages-user-and-friend'
    ];
}
