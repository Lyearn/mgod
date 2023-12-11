import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';

function HomepageHeader() {
  return (
    <header className={clsx('hero hero--secondary', styles.heroBanner)}>
      <div className={clsx('container', styles.heroContainer)}>
        <Heading as="h1" className="hero__title">
          Empower Your Go Applications with Mgod
        </Heading>
        <p className="hero__subtitle">
          Transform your MongoDB interactions in Go effortlessly with mgod. Simplify database operations, enhance type safety, and build robust applications!
        </p>
        <div className={styles.buttons}>
          <Link
            className="button button--primary button--lg"
            to="/docs/intro">
            Get Started
          </Link>
        </div>
      </div>
    </header>
  );
}

function HomepageFooter() {
  return (
    <footer className={clsx('hero hero--secondary', styles.ctaBanner)}>
      <div className={clsx('container', styles.ctaContainer)}>
        <Heading as="h2" className="hero__title">
        Ready to Simplify Your MongoDB Interactions?
        </Heading>
        <div className={styles.buttons}>
          <Link
            className="button button--primary button--lg"
            to="/docs/introduction/about">
            Get Started
          </Link>
        </div>
      </div>
    </footer>
  );
}

export default function Home(): JSX.Element {
  return (
    <Layout
      title="Home">
      <HomepageHeader />
      <main>
        <HomepageFeatures />
      </main>
      <HomepageFooter />
    </Layout>
  );
}
