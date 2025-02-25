const AuthManager = (function() {
    const tokenStore = new WeakMap();
    const context = Object.freeze({});

    return Object.freeze({
        setToken(token) {
            if (typeof token !== 'string' && token !== null) {
                throw new TypeError('Token must be a string or null');
            }
            tokenStore.set(context, token);
        },

        getToken() {
            return tokenStore.get(context);
        },

        clearToken() {
            tokenStore.delete(context);
        },

        async validateToken() {
            const token = this.getToken();
            if (!token) return false;

            try {
                const response = await fetch('/api/admin/validate', {
                    method: 'GET',
                    headers: {
                        'Authorization': token,
                        'Content-Type': 'application/json'
                    }
                });
                return response.ok;
            } catch {
                return false;
            }
        }
    });
})();

export const auth = AuthManager;
