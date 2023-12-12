import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

function GetStartedButton() {
  return (
    <Link
      className="px-[3.2rem] py-[1.2rem] rounded-full heading-xs text-white bg-black hover:text-white focus-visible:text-white"
      to="/docs/introduction/about">
      Get Started
    </Link>
  );
}

function HomepageHeader() {
  return (
    <header className="py-[12rem] text-center">
      <Heading as="h1" className="pb-[3.2rem] max-w-[82.5rem] mx-auto text-text-primary heading-xxxl" style={{
        fontSize: '6.8rem',
        fontStyle: 'normal',
        fontWeight: 600,
        lineHeight: 'normal',
        letterSpacing: '-0.136rem',
      }}>
        Empower Your Go Applications with Mgod
      </Heading>
      <p className="pb-[5.6rem] max-w-[65.4rem] mx-auto text-text-secondary body-long-01" style={{
        fontSize: '1.8rem',
        fontStyle: 'normal',
        fontWeight: 400,
        lineHeight: '3.2rem',
      }}>
        Transform your MongoDB interactions in Go effortlessly with mgod. Simplify database operations, enhance type safety, and build robust applications!
      </p>
      <GetStartedButton />
    </header>
  );
}

function HomepageFooter() {
  return (
    <footer className="py-[20rem] text-center">
      <Heading as="h2" className="pb-[5.6rem] max-w-[59.6rem] mx-auto text-text-primary heading-xxl" style={{
        fontSize: '4.8rem',
        fontStyle: 'normal',
        fontWeight: 600,
        lineHeight: '140%',
        letterSpacing: '-0.096rem',
      }}>
        Ready to Simplify Your MongoDB Interactions?
      </Heading>
      <GetStartedButton />
    </footer>
  );
}

export default function Home(): JSX.Element {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title="mgod - Empower your Go applications"
      description={siteConfig.tagline}>
      <HomepageHeader />
      <main>
        <HomepageFeatures />
      </main>
      <HomepageFooter />
    </Layout>
  );
}
