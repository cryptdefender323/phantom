import AsciinemaPlayer from "@/components/asciinema";
import { PhantomsIcon } from "@/components/icons/phantoms";
import TutorialCard from "@/components/tutorial-card";
import { Themes } from "@/util/themes";
import { faDownload, faExternalLink } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Button, Card, CardBody, CardHeader, Divider, Link } from "@heroui/react";
import { useTheme } from "next-themes";
import { useRouter } from "next/router";
import React from "react";

export default function Home() {
  const router = useRouter();
  const { resolvedTheme } = useTheme();
  const [mounted, setMounted] = React.useState(false);

  React.useEffect(() => {
    setMounted(true);
  }, []);

  const isDarkTheme = React.useMemo(() => {
    if (!mounted) {
      return true;
    }
    return (resolvedTheme || Themes.DARK) !== Themes.LIGHT;
  }, [mounted, resolvedTheme]);

  return (
    <div className="mt-6 flex flex-col gap-6 px-4 sm:px-6 lg:grid lg:grid-cols-12 lg:px-12">
      <div className="lg:col-span-6">
        <div className="w-full overflow-hidden rounded-xl border border-default-200 bg-content1 shadow-sm dark:border-default-100/60">
          <div className="w-full overflow-x-auto lg:overflow-visible">
            <AsciinemaPlayer
              src="/asciinema/intro.cast"
              rows="18"
              cols="75"
              idleTimeLimit={60}
              preload={true}
              autoPlay={true}
              loop={true}
            />
          </div>
        </div>
      </div>
      <div className="lg:col-span-6 lg:ml-2">
        <Card className="mx-auto max-w-3xl lg:mx-0">
          <CardHeader>
            <div className="flex items-center">
              <PhantomsIcon className="mr-2" height={28} />
              <span className="text-2xl">Phantom Command &amp; Control</span>
            </div>
          </CardHeader>
          <Divider />
          <CardBody>
            <p className={isDarkTheme ? "prose dark:prose-invert" : "prose prose-slate"}>
              Phantom is a powerful command and control (C2) framework designed
              to provide advanced capabilities for covertly managing and
              controlling remote systems. With Phantom, security professionals,
              red teams, and penetration testers can easily establish a secure
              and reliable communication channel over Mutual TLS, HTTP(S), DNS,
              or Wireguard with target machines. Enabling them to execute
              commands, gather information, and perform various
              post-exploitation activities. The framework offers a user-friendly
              console interface, extensive functionality, and support for
              multiple operating systems as well as multiple CPU architectures,
              making it an indispensable tool for conducting comprehensive
              offensive security operations.
            </p>
            <div className="mt-4 flex w-full gap-3">
              <Button
                className="flex-1"
                color="primary"
                variant="shadow"
                as={Link}
                href="https://github.com/cryptdefender323/phantom/releases/latest"
                target="_blank"
                rel="noopener noreferrer"
                startContent={<FontAwesomeIcon icon={faDownload} />}
              >
                Download Latest Release
              </Button>
              <Button
                className="flex-1"
                color="secondary"
                variant="ghost"
                as={Link}
                href="https://github.com/phantomarmory"
                target="_blank"
                rel="noopener noreferrer"
                endContent={<FontAwesomeIcon icon={faExternalLink} />}
              >
                Visit the Armory
              </Button>
            </div>
          </CardBody>
        </Card>
      </div>

      <div className="col-span-12 mt-8">
        <Divider />
      </div>

      <div className="col-span-12 mt-8">
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-9">
          <div className="sm:col-span-1 lg:col-span-3">
            <TutorialCard
              name="Getting Started"
              description="A quick start guide to get you up and running"
              asciiCast="/asciinema/install-1.cast"
              cols="133"
              rows="32"
              idleTimeLimit={1}
              showButton={true}
              buttonText="Read Docs"
              onPress={() => {
                router.push({
                  pathname: "/docs",
                  query: { name: "Getting Started" },
                });
              }}
            />
          </div>

          <div className="sm:col-span-1 lg:col-span-3">
            <TutorialCard
              name="Compile From Source"
              description="How to compile Phantom from source"
              asciiCast="/asciinema/compile-from-source.cast"
              cols="133"
              rows="32"
              idleTimeLimit={1}
              showButton={true}
              buttonText="Read Docs"
              onPress={() => {
                router.push({
                  pathname: "/docs",
                  query: { name: "Compile from Source" },
                });
              }}
            />
          </div>
        </div>
      </div>

      <div className="col-span-12 mb-8"></div>
    </div>
  );
}
