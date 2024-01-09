import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

function GetStartedButton() {
  return (
    <Link
      className="block w-[14rem] py-[1.2rem] mx-auto rounded-full heading-xs text-white bg-black hover:text-white focus-visible:text-white dark:bg-[#222]"
      to="/docs/about">
      Get Started
    </Link>
  );
}

function HomepageHeader() {
  return (
    <header className="pt-[6.4rem] sm:pt-[12rem] pb-[12rem] px-[2rem] text-center">
      <div className="max-w-[82.5rem] mx-auto">
        <Heading
          as="h1"
          className="mb-[3.2rem] mx-auto sm:text-[6.8rem] sm:leading-[normal] sm:tracking-[-0.136rem] font-[600] text-[4.8rem] leading-[120%] tracking-[-0.096rem] text-text-primary dark:text-text-white">
          Empower your Go applications with mgod
        </Heading>
        <p className="mb-[5.6rem] max-w-[clamp(34.2rem,90%,65.4rem)] mx-auto text-[1.8rem] font-[400] leading-[3.2rem] text-text-secondary">
          Transform your MongoDB interactions in Go effortlessly with mgod. Simplify database
          operations, enhance type safety, and build robust applications!
        </p>
      </div>
      <GetStartedButton />
    </header>
  );
}

function HomepageFooter() {
  return (
    <footer className="py-[12rem] sm:py-[20rem] px-[2rem] text-center">
      <Heading
        as="h2"
        className="mb-[5.6rem] max-w-[59.6rem] mx-auto text-text-primary dark:text-text-white text-[3.2rem] font-[600] leading-[120%] tracking-[-0.064] sm:text-[4.8rem] sm:leading-[140%] sm:tracking-[-0.096rem]">
        Ready to simplify your MongoDB interactions?
      </Heading>
      <GetStartedButton />
    </footer>
  );
}

function Home() {
  const { siteConfig } = useDocusaurusContext();

  return (
    <Layout title="mgod - Empower your Go applications" description={siteConfig.tagline}>
      <HomepageHeader />
      <main>
        <HomepageFeatures />
      </main>
      <HomepageFooter />
    </Layout>
  );
}

export default Home;
