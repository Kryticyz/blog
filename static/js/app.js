const { useState, useEffect } = React;

function App() {
	const [posts, setPosts] = useState([]);

	useEffect(() => {
		fetch('/api/posts')
			.then(res => res.json())
			.then(data => setPosts(data));
	}, []);

	return React.createElement('div', { className: 'container' },
		React.createElement('h1', { className: 'title' }, 'Blog'),
		React.createElement('div', { className: 'posts' },
			posts.map(post =>
				React.createElement('a', {
					key: post.slug,
					href: `/post/${post.slug}`,
					className: 'post-card'
				},
					post.image && React.createElement('img', { src: `/${post.image}`, alt: post.title }),
					React.createElement('h2', null, post.title),
					React.createElement('time', null, new Date(post.date).toLocaleDateString('en-US', {
						year: 'numeric',
						month: 'long',
						day: 'numeric'
					}))
				)
			)
		)
	);
}

ReactDOM.render(React.createElement(App), document.getElementById('root'));
